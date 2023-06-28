import { LoaderFunction, redirect } from "react-router-dom";
import { Session, SessionOverview } from "./data";

export interface LoaderResult {
  sessionOverview: SessionOverview;
  session?: Session;
}

export const SessionOverviewLoader: LoaderFunction = async (): Promise<LoaderResult | Response> => {
  const sessionOverview = await loadSessionOverview();
  if (sessionOverview.sessionSummaries.length !== 0) {
    return redirect(`/session/${sessionOverview.sessionSummaries[0].id}`);
  }
  return {sessionOverview};
};

export const SessionLoader: LoaderFunction = async ({params}): Promise<LoaderResult | Response> => {
  const sessionOverview = await loadSessionOverview();

  const sessionId = params.sessionId;
  if (sessionId == null) {
    return redirect('/');
  }

  const validSessionIds = sessionOverview.sessionSummaries.map(s => s.id);
  if ( ! validSessionIds.includes(sessionId)) {
    return redirect('/');
  }

  const session = await loadSession(sessionId);

  return { sessionOverview, session };
};

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
      title: '',
      prompt: '',
      responses: []
    };
  }
  return {
    id: '2222',
    title: '',
    prompt: '',
    responses: []
};
}
