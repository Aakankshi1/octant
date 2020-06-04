import { withKnobs } from '@storybook/addon-knobs';
import { addDecorator, moduleMetadata } from '@storybook/angular';
import { SharedModule } from '../src/app/modules/shared/shared.module';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';
import { setConsoleOptions } from '@storybook/addon-console';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {APP_INITIALIZER} from "@angular/core";
import {InitService} from "../src/app/modules/shared/services/init/init.service";

// enable knobs
addDecorator(withKnobs);

setConsoleOptions({
  panelExclude: [
    /Angular is running in the development mode/,
    /Ignored an update to unaccepted module/,
  ],
});

addDecorator(
  moduleMetadata({
    imports: [
      BrowserAnimationsModule,
      SharedModule,
      MarkdownModule.forRoot({
        markedOptions: {
          provide: MarkedOptions,
          useValue: {
            gfm: true,
            tables: true,
            breaks: true,
            pedantic: false,
            sanitize: false,
            smartLists: true,
            smartypants: false,
          },
        },
      }),
    ],
    providers: [
      InitService,
      {
        provide: APP_INITIALIZER,
        useFactory: (initService) => () => initService.init(),
        deps: [InitService],
        multi: true
      },
    ]
  })
);
