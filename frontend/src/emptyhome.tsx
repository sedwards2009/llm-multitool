import { useEffect } from "react";
import { SessionOverview } from "./data";
import { navigate } from "raviger";

export interface Props {
  sessionOverview: SessionOverview;
}

export function EmptyHome({ sessionOverview }: Props): JSX.Element {
  useEffect(() => {
    if (sessionOverview.sessionSummaries.length === 0) {
      return;
    }
    navigate(`/session/${sessionOverview.sessionSummaries[0].id}`);
  }, [sessionOverview]);

  return <div>Empty Home</div>;
}
