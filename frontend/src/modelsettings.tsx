import { } from "react";
import { ModelOverview, TemplateOverview } from "./data";

export interface Props {
  modelOverview: ModelOverview;
  templateOverview: TemplateOverview;
  selectedModelId: string | null;
  setSelectedModelId: (id: string) => void;
  selectedTemplateId: string | null;
  setSelectedTemplateId: (id: string) => void;
}

export function ModelSettings({modelOverview, templateOverview, selectedModelId, setSelectedModelId,
    selectedTemplateId, setSelectedTemplateId}: Props): JSX.Element  {

  return <div className="gui-layout cols-1-3">
    <div>
      <i className="fa fa-robot"></i>&nbsp;&nbsp;Model:
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
    <div>
      <i className="fa fa-hammer"></i>&nbsp;&nbsp;Task:
    </div>
    <div>
      <select
        value={selectedTemplateId == null ? undefined : selectedTemplateId}
        onChange={(e) => setSelectedTemplateId(e.target.value)}>
          <option key={""} value="">(none)</option>
        {
          templateOverview.templates.map(t => <option key={t.id} value={t.id}>{ t.name }</option>)
        }
      </select>
    </div>
  </div>;
}