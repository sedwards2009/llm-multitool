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
    return await response.json();
  } catch (error) {
    console.error("Could not parse JSON", error);
    return null;
  }
}
