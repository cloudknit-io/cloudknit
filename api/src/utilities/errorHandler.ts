import { BadRequestException } from "@nestjs/common";
import { SqlErrorCodes } from "src/types";

export function handleSqlErrors(err: any, dupEntryMsg: string = 'entry already exists') {
  if (!err || !err.code) {
    return;
  }
  
  if (err.code === SqlErrorCodes.DUP_ENTRY) {
    throw new BadRequestException(dupEntryMsg);
  }

  if (err.code === SqlErrorCodes.NO_DEFAULT) {
    throw new BadRequestException(err.sqlMessage);
  }
}
