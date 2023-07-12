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


export enum SessionMonitorState {
  IDLE,
  CONNECTING,
  CONNECTED,
  WAITING_TO_RECONNECT,
}

const DEFAULT_RECONNECT_DELAY_MS = 100;

export class SessionMonitor {
  #sessionId = "";
  #socket: WebSocket | null = null;
  #state = SessionMonitorState.IDLE;
  #onChange: ((message: string) => void)  | null = null;
  #onStateChange: ((statue: SessionMonitorState) => void)  | null = null;
  #reconnectDelayMs = DEFAULT_RECONNECT_DELAY_MS;

  constructor(sessionId: string, onChange: (message: string) => void,
      onStateChange: (state: SessionMonitorState) => void) {
console.log(`new SessionMonitor`);
    this.#sessionId = sessionId;
    this.#onChange = onChange;
    this.#onStateChange = onStateChange;
  }

  state(): SessionMonitorState {
    return this.#state;
  }

  #setState(state: SessionMonitorState): void {
    this.#state = state;
    if (this.#onStateChange != null) {
      console.log(`SessionMonitor ${state}`);
      this.#onStateChange(state);
    }
  }

  start(): void {
    this.#connect();
  }

  #connect(): void {
    console.log(`Connecting`);
    this.#setState(SessionMonitorState.CONNECTING);
    this.#socket = new WebSocket(`${WEBSOCKET_SERVER_BASE_URL}/session/${this.#sessionId}/changes`)
    this.#socket.addEventListener("message", (event) => {
      console.log(`Received Message: ${event.data}`);
      if (this.#onChange != null) {
        this.#onChange(event.data);
      }
    });
    this.#socket.addEventListener("open", () => {
      console.log(`Websocket open for sessionId ${this.#sessionId}`);
      this.#setState(SessionMonitorState.CONNECTED);
      this.#reconnectDelayMs = DEFAULT_RECONNECT_DELAY_MS;
    });
    this.#socket.addEventListener("close", () => {
      console.log(`Websocket closed for sessionId ${this.#sessionId}`);
      if (this.#state === SessionMonitorState.IDLE) {
        return;
      }
      this.#reconnect();
    });
    this.#socket.addEventListener("error", (e) => {
      console.log(`Websocket error for sessionId ${this.#sessionId}`, e);
      if (this.#socket !== null) {
        this.#socket.close();
      }
    });
  }

  #reconnect(): void {
    this.#setState(SessionMonitorState.WAITING_TO_RECONNECT);
    console.log(`Reconnecting ${this.#reconnectDelayMs}ms`);
    setTimeout(() => {
      if (this.#state === SessionMonitorState.IDLE) {
        return;
      }
      this.#connect();
    }, this.#reconnectDelayMs);
    this.#reconnectDelayMs = Math.min(this.#reconnectDelayMs * 2, 5000);
  }

  stop(): void {
    this.#state = SessionMonitorState.IDLE;
    if (this.#socket == null) {
      return;
    }
    this.#socket.close();
    this.#socket = null;
  }
}