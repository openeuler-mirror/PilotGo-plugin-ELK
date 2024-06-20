/// <reference types="vite/client" />

declare module "*.vue";
interface Window {
    remount: any;
    unmount: any;
    readonly '__MICRO_APP_BASE_ROUTE__': string
}
