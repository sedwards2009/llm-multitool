import classNames from "classnames";
import { SessionOverview, SessionSummary } from "./data";
import { navigate } from "raviger";

export interface Props {
  sessionOverview: SessionOverview;
  selectedSessionId: string | null;
}

export function SessionOverviewList({sessionOverview, selectedSessionId}: Props): JSX.Element {
  const sessionSummaries = [...sessionOverview.sessionSummaries];
  const cmp = (a: SessionSummary, b: SessionSummary) => {
    if (a.creationTimestamp < b.creationTimestamp) {
      return -1;
    }
    return a.creationTimestamp === b.creationTimestamp ? 0 : 1;
  }
  sessionSummaries.sort(cmp);

  return <ul className="tabs">
    {
      sessionSummaries.map(s => {
        const onClick = () => {
          navigate(`/session/${s.id}`);
        };
        return <li
          key={s.id}
          className={classNames({"tab": true, active: s.id === selectedSessionId})}
          onClick={onClick}
          >
            {s.title}
          </li>;
      })
    }
  </ul>;
}
