import { HttpClient, HttpResponse } from '@angular/common/http';
import { ListTopics } from './topicService';
import { tap, map, mergeMap, switchMap } from 'rxjs/operators';
import { Observable, of, interval, forkJoin } from 'rxjs';
import { ajax } from 'rxjs/ajax';


export interface TopicInfo {
    Data: string[]
    Topic: string
}

export interface ClientTopicInfo {
    data: string
    name: string
}

export type ClientData = Array<TopicInfo>

export class ClientService {
  private host = "https://develop.valenoq.com/sandbox/khanty-app/client"

  constructor(private http: HttpClient) {

  }

  public getState(port: string, topics: string[]) :  Observable<ClientTopicInfo> {
      const url = this.host + "/" + port
      var urls = {}
      for (const topicName of topics) {
          urls[topicName] = ajax.getJSON(url + "/" + topicName)
      }

      return interval(1000).pipe(switchMap(() => forkJoin(urls)))
  }
}
