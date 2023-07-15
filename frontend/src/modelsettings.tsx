import { } from "react";
import { ModelOverview } from "./data";

export interface Props {
  modelOverview: ModelOverview;
  selectedModelId: string | null;
  setSelectedModelId: (id: string) => void;
}

export function ModelSettings({modelOverview, selectedModelId, setSelectedModelId}: Props): JSX.Element  {
  return <div className="gui-layout cols-1-3">
    <div>
      <i className="fa fa-brain"></i> Model:
    </div>
    <div>
      <select
        value={selectedModelId == null ? undefined : selectedModelId}
        onChange={(e) => setSelectedModelId(e.target.value)}>
          <option key={""} value="">(none)</option>
        {
          modelOverview.models.map(m => <option key={m.id} value={m.id}>{ m.name }</option>)
        }
      </select>
    </div>
  </div>;
}