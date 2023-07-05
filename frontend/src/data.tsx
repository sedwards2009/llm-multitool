
export interface SessionOverview {
  sessionSummaries: SessionSummary[];
}

export interface SessionSummary {
  id: string;
  title: string;
  creationTimestamp: string;
}

export interface Root {
  sessions: Session[];
}

export interface Session {
  id: string;
  creationTimestamp: string;
  title: string;
  prompt: string;
  responses: Response[];
}

export type ResponseStatus = "done" | "pending" | "running";

export interface Response {
  id: string;
  creationTimestamp: string;
  status: ResponseStatus;
  prompt: string;
  text: string;
}
