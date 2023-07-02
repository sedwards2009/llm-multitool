import * as _ from "lodash";
import { Session, SessionOverview } from "./data";

export interface LoaderResult {
  sessionOverview: SessionOverview;
  session?: Session;
}

const SERVER_BASE_URL = "http://localhost:8080";

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
