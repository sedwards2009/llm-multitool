import { useEffect } from "react";
import { SessionOverview } from "./data";
import { navigate } from "raviger";
import { NewSessionButton } from "./newsessionbutton";

export interface Props {
  sessionOverview: SessionOverview;
  onSessionChange: ()=> void;
}

export function EmptyHome({ sessionOverview, onSessionChange }: Props): JSX.Element {
  useEffect(() => {
    if (sessionOverview.sessionSummaries.length === 0) {
      return;
    }
    navigate(`/session/${sessionOverview.sessionSummaries[0].id}`);
  }, [sessionOverview]);

  return <div>
    <NewSessionButton onSessionChange={onSessionChange} />
  </div>;
}
