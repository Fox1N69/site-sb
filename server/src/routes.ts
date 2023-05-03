import {Express, Request, Response} from 'express';

export default function (app: Express) {
  app.get('/', (res: Response, req: Request) => res.sendStatus(200));
}

