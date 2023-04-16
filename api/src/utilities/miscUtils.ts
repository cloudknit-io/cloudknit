export class MiscUtils {
  static async wait(ms: number) {
    return new Promise((resolve, reject) => setTimeout(resolve, ms));
  }
}
