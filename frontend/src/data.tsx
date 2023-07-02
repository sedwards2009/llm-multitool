
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
  title: string;
  prompt: string;
  responses: Response[];
}

export interface Response {
  text: string;
}
