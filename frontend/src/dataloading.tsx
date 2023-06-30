import { Session, SessionOverview } from "./data";

export interface LoaderResult {
  sessionOverview: SessionOverview;
  session?: Session;
}

export async function loadSessionOverview(): Promise<SessionOverview> {
  const sessionOverview: SessionOverview = {
    sessionSummaries: [
      {
        id: '1111',
        title: 'Qt event questions'
      },
      {
        id: '2222',
        title: 'Simple React component'
      },
    ]
  };
  return sessionOverview;
}

export async function loadSession(sessionId: string): Promise<Session> {
  if (sessionId === '1111') {
    return {
      id: '1111',
      title: 'Qt event questions',
      prompt: 'Which Qt event is for a window gaining focus?',
      responses: []
    };
  }
  return {
    id: '2222',
    title: 'Simple React component',
    prompt: 'Write out a simple React component.',
    responses: []
  };
}
