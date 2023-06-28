import { SessionOverviewList } from "./sessionoverviewlist";
import { useLoaderData } from "react-router-dom";
import { SessionEditor } from "./sessioneditor";
import { LoaderResult } from "./dataloading";

export function App() {
  const loaderData = useLoaderData() as LoaderResult;
  return (
    <>
      <h1>LLM Workbench</h1>
      <div className="top-layout">
        <div className="session-list">
          <SessionOverviewList
            sessionOverview={loaderData.sessionOverview}
            selectedSessionId={loaderData.session?.id ?? null}
          />
          <button>+ New Session</button>
        </div>
        <div className="session-tab">
        {loaderData.session &&
          <SessionEditor
            session={loaderData.session}
          />}
        </div>
      </div>
    </>
  )
}
