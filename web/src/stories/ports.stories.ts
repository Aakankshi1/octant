import {PortsView} from "../app/modules/shared/models/content";
import {storiesOf} from "@storybook/angular";
import {PortsComponent} from "../app/modules/shared/components/presentation/ports/ports.component";
import {object} from "@storybook/addon-knobs";

const portsView1: PortsView = {
  config: {
    ports: [
      {
        metadata: {
          type: 'port',
        },
        config: {
          port: 80,
          protocol: 'TCP',
          state: {
            isForwardable: true,
          },
          buttonGroup: {
            metadata: {
              type: 'buttonGroup',
            },
            config: {
              buttons: [
                {
                  name: 'Start port forward',
                  payload: {
                    action: 'overview/startPortForward',
                    apiVersion: 'v1',
                    kind: 'Pod',
                    name: 'httpbin-db6d74d85-nltjq',
                    namespace: 'default',
                    port: 80,
                  },
                },
              ],
            },
          },
        },
      },
    ],
  },
  metadata: {
    type: 'ports',
  },
};

storiesOf('Components', module).add('Ports', () => ({
  props: {
    view: object('View', portsView1)
  },
  component: PortsComponent,
}));
