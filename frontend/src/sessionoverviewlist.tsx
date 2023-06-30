import classNames from "classnames";
import { SessionOverview } from "./data";
import { navigate } from "raviger";

export interface Props {
  sessionOverview: SessionOverview;
  selectedSessionId: string | null;
}

export function SessionOverviewList({sessionOverview, selectedSessionId}: Props): JSX.Element {
  return <ul className="tabs">
    {
      sessionOverview.sessionSummaries.map(s => {
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
