import React, { useMemo } from 'react';
import { SceneApp, SceneAppPage } from '@grafana/scenes';
import { ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/utils.routing';
import { getDeliveryServiceScene } from './scene';

const getScene = () =>
  new SceneApp({
    pages: [
      new SceneAppPage({
        title: 'Delivery Services',
        url: prefixRoute(`${ROUTES.DeliveryService}`),
        hideFromBreadcrumbs: true,
        getScene: getDeliveryServiceScene,
      }),
    ],
  });

export const DeliveryServicePage = () => {
  const scene = useMemo(() => getScene(), []);

  return <scene.Component model={scene} />;
};
