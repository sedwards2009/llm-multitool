import { useRoutes } from "raviger";
import { SessionOverview } from "./data";
import { EmptyHome } from "./emptyhome";
import { Home } from "./home";

export interface Props {
  sessionOverview: SessionOverview;
}

export function MainApp({ sessionOverview }: Props): JSX.Element {
  return (
    <>
      <h1>LLM Workbench</h1>
      {
        useRoutes(
          {
            '/': () => <EmptyHome sessionOverview={sessionOverview} />,
            '/session/:sessionId': ({ sessionId }: { sessionId: any }) => {
              return <Home
                sessionOverview={sessionOverview}
                sessionId={sessionId}
              />;
            }
          }
        )
      }
    </>
  );
}
