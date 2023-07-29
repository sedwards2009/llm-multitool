import { navigate, useQueryParams } from "raviger";
import { ModelOverview } from "./data";


export interface Props {
  modelOverview: ModelOverview;
}

export function SettingsPage({ modelOverview }: Props): JSX.Element {
  const [{from}, _] = useQueryParams();
  const onBackClicked = () => {
    navigate(from);
  };

  return (
    <>
      <h2>Settings</h2>

      <h3>Models</h3>
      <ul>
      {
        modelOverview.models.map(m => {
          return <li key={m.id}>{m.name}</li>
        })
      }
      </ul>
      <button className="small primary" onClick={onBackClicked}>
        <i className="fas fa-arrow-left"></i>
        {" Back"}
      </button>
    </>
  );
}
