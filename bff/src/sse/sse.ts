import helper from "../utils/helper";
import * as express from "express";
import { createClient } from "redis";
import EventEmitter = require("events");

const event = new EventEmitter();

async function startRedis() {
  console.log('Starting Redis');
  const client = createClient({
    url: "http://cloudknit-redis-master.zlifecycle-system.svc.cluster.local:6379"
  });

  client.on("error", (err) => console.log("Redis Client Error", err));

  await client.connect();
  client.SUBSCRIBE("test-channel", (message: any, channel: any) => {
    console.log(JSON.parse(message), channel);
    event.emit("stream", JSON.parse(message));
  });
}

function eventsHandler(request: any, response: any, next) {
  const reqUser = helper.userFromReq(request);
  if (!reqUser) {
    response.status(400).send({
      unauthorized: "User is not authorized",
    });
    return;
  }
  const headers = {
    "Content-Type": "text/event-stream",
    Connection: "keep-alive",
    "Cache-Control": "no-cache",
  };

  response.writeHead(200, headers);

  response.write(`data: ${JSON.stringify({ d: new Date() })}\n\n`);

  event.on("stream", (data: string) => {
    response.write(`data: ${data}\n\n`);
  });

  request.on("close", () => {
    console.log(`Connection closed`);
  });
}

export function setUpSSE(app: express.Express) {
  startRedis().then(() => {});
  app.get("/stream", eventsHandler);
}
