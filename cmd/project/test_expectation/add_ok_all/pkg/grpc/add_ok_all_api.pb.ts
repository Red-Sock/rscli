/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "../fetch.pb";
import * as GoogleProtobufTimestamp from "../google/protobuf/timestamp.pb";


export type PingRequest = {
  clientTimestamp?: GoogleProtobufTimestamp.Timestamp;
};

export type PingResponse = {
  took?: number;
};

export class addOkAllAPI {
  static Version(this:void, req: PingRequest, initReq?: fm.InitReq): Promise<PingResponse> {
    return fm.fetchRequest<PingResponse>(`/v1/version`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}