import classNames from "classnames";
import { SessionOverview, SessionSummary } from "./data";
import { navigate } from "raviger";
import groupBy from "just-group-by";


export interface Props {
  sessionOverview: SessionOverview;
  selectedSessionId: string | null;
}

export function SessionOverviewList({sessionOverview, selectedSessionId}: Props): JSX.Element {
  const groups = groupBy(sessionOverview.sessionSummaries, (summary: SessionSummary): string => {
    return summary.creationTimestamp.substring(0, 10);
  });
  console.log(groups);

  const groupNames = Object.getOwnPropertyNames(groups);
  groupNames.sort().reverse();
  return <>
  {
    groupNames.map((name) => formatGroup(groups[name], selectedSessionId))
  }
  </>;
}

function formatGroup(sessionSummaries: SessionSummary[], selectedSessionId: string | null): JSX.Element {
  const sortedSessionSummaries = [...sessionSummaries];
  const cmp = (a: SessionSummary, b: SessionSummary) => {
    if (a.creationTimestamp > b.creationTimestamp) {
      return -1;
    }
    return a.creationTimestamp === b.creationTimestamp ? 0 : 1;
  }
  sortedSessionSummaries.sort(cmp);

  const dateFormat = new Intl.DateTimeFormat("default", { dateStyle: "medium" });

  const now = new Date();
  const groupDay = new Date(Date.parse(sortedSessionSummaries[0].creationTimestamp));
  const isToday = (now.getDay() === groupDay.getDay() &&
    now.getMonth() === groupDay.getMonth() &&
    now.getFullYear() === groupDay.getFullYear());

  return <>
    {
      !isToday && <p className="minor">
        {
          dateFormat.format(Date.parse(sortedSessionSummaries[0].creationTimestamp))
        }
      </p>
    }
    <ul className="tabs">
      {
        sortedSessionSummaries.map(s => {
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
    </ul>
  </>;
}
