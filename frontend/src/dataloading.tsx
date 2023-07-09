import * as _ from "lodash";
import { Session, SessionOverview } from "./data";

export interface LoaderResult {
  sessionOverview: SessionOverview;
  session?: Session;
}

const SERVER_BASE_URL = "http://localhost:8080";
const WEBSOCKET_SERVER_BASE_URL = "ws://localhost:8080";

export async function loadSessionOverview(): Promise<SessionOverview> {
  const response = await fetch(`${SERVER_BASE_URL}/session`);
  try {
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return {
      sessionSummaries: []
    };
  }
}

export async function loadSession(sessionId: string): Promise<Session | null> {
  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}`);
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

export async function newSession(): Promise<Session | null> {
  const response = await fetch(`${SERVER_BASE_URL}/session`, {method: "POST"});
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

let saveSessionPromptQueue = new Map<string, string>();

export function setSessionPrompt(session: Session, prompt: string): Session {
  saveSessionPromptQueue.set(session.id, prompt);
  flushSessionPromptQueue();
  return {...session, prompt};
}

async function processSessionPromptQueue(): Promise<void> {
  const workingSessionPromptQueue = saveSessionPromptQueue;
  saveSessionPromptQueue = new Map<string, string>();
  for (const [sessionId, prompt] of workingSessionPromptQueue.entries()) {
    await fetch(`${SERVER_BASE_URL}/session/${sessionId}/prompt`,
      {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({prompt})
    });
  }
}

const flushSessionPromptQueue = _.throttle(() => {
  (async () => {
    processSessionPromptQueue();
  })();
}, 1000);

async function flushQueues(): Promise<void> {
  await processSessionPromptQueue();
}

export async function newResponse(session: Session): Promise<Response | null> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${session.id}/response`, {method: "POST"});
  try {
    if (response.ok) {
      return await response.json();
    }
  } catch (error) {
    console.error("Could not parse JSON", error);
  }
  return null;
}

export async function deleteResponse(sessionId: string, responseId: string): Promise<boolean> {
  await flushQueues();

  const response = await fetch(`${SERVER_BASE_URL}/session/${sessionId}/response/${responseId}`, {method: "DELETE"});
  return response.ok;
}

export class SessionMonitor {
  #sessionId = "";
  #socket: WebSocket | null = null;
  #callback: ((message: string) => void)  | null = null;

  constructor(sessionId: string, callback: (message: string) => void) {
    this.#sessionId = sessionId;
    this.#callback = callback;
  }

  start(): void {
    this.#socket = new WebSocket(`${WEBSOCKET_SERVER_BASE_URL}/session/${this.#sessionId}/changes`)
    this.#socket.addEventListener("message", (event) => {
      console.log(`Received Message: ${event.data}`);
      if (this.#callback != null) {
        this.#callback(event.data);
      }
    });
    this.#socket.addEventListener("open", () => {
      console.log(`Websocket open for sessionId ${this.#sessionId}`);
    });
    this.#socket.addEventListener("close", () => {
      console.log(`Websocket closed for sessionId ${this.#sessionId}`);
    });
    this.#socket.addEventListener("error", (e) => {
      console.log(`Websocket error for sessionId ${this.#sessionId}`, e);
    });
  }

  stop(): void {
    if (this.#socket == null) {
      return;
    }
    this.#socket.close();
    this.#socket = null;
  }
}