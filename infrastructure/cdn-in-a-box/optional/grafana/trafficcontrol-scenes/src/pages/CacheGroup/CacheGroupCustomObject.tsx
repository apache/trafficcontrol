import React from 'react';

import { SceneComponentProps, SceneObjectBase, SceneObjectState } from '@grafana/scenes';
import { InlineField, Input } from '@grafana/ui';

interface CacheGroupState extends SceneObjectState {
  name: string;
}

function CacheGroupInputRenderer({ model }: SceneComponentProps<CacheGroupCustomObject>) {
  const { name } = model.useState();

  return (
    <InlineField label="cachegroup" style={{ margin: '0' }}>
      <Input
        prefix=""
        defaultValue={name}
        width={20}
        type="string"
        onBlur={(evt) => {
          model.onValueChange(evt.currentTarget.value);
        }}
      />
    </InlineField>
  );
}

export class CacheGroupCustomObject extends SceneObjectBase<CacheGroupState> {
  public static Component = CacheGroupInputRenderer;

  onValueChange = (value: string) => {
    this.setState({ name: value });
  };
}
