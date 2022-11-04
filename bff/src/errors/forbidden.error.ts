export default class ForbiddenError extends Error {
  code=null;

  constructor(message: string, code: number = 403) {
    super(message);
    this.code = code;
  }
}
