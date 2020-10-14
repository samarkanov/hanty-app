import { HttpClient, HttpResponse } from '@angular/common/http';
import { tap, map, mergeMap } from 'rxjs/operators';
import { Observable, interval } from 'rxjs';


export interface ConfigItem {
    host: string
    port: string
    name?: string
    message?: string
    color?: string
    messagesSubscribe?: boolean
    colorSubscribe?: boolean
}

export type Config = Record<string, ConfigItem>

export class ConfigService {
  private host = "https://develop.valenoq.com/sandbox/khanty-app/config"
  private url: string

  constructor(private http: HttpClient) {
      this.url = this.host
  }

  public getConfig() : Observable<Config> {
      return this.http.get<Config>(this.url)
  }
}
