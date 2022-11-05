declare namespace NodeJS {
  export interface ProcessEnv {
    MY_SQL_HOST: string;
    MY_SQL_PORT: string;
    MY_SQL_USERNAME: string;
    MY_SQL_PASSWORD: string;
    MY_SQL_DATABASE: string;
    APP_PORT: string;
  }
}
