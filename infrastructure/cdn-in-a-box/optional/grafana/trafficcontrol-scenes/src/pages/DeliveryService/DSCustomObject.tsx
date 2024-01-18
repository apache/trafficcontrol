import React from 'react';

import { SceneComponentProps, SceneObjectBase, SceneObjectState } from '@grafana/scenes';
import { InlineField, Input } from '@grafana/ui';

interface DeliveryServiceState extends SceneObjectState {
  name: string;
}

function Renderer({ model }: SceneComponentProps<DeliveryServiceCustomObject>) {
  const { name } = model.useState();

  return (
    <InlineField label="deliveryservice" style={{ margin: '0' }}>
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

export class DeliveryServiceCustomObject extends SceneObjectBase<DeliveryServiceState> {
  public static Component = Renderer;

  onValueChange = (value: string) => {
    this.setState({ name: value });
  };
}
