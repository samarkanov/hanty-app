import { HttpClient, HttpResponse } from '@angular/common/http';
import { tap, mergeMap, map, switchMap } from 'rxjs/operators';
import { Observable, interval, forkJoin } from 'rxjs';
import { ajax } from 'rxjs/ajax';

export interface ListTopics {
    topics: string[]
}

export interface StateTopic {
    states: string[]
    topic: string
}

export class TopicService {
  private host = "https://develop.valenoq.com/sandbox/khanty-app/topic"
  private url: string

  constructor(private http: HttpClient) {
      this.url = this.host
  }

  public listTopics() : Observable<ListTopics> {
      return this.http.get<ListTopics>(this.url);
  }

  public statesForTopic(topic: string) : Observable<StateTopic> {
      return this.http.get<StateTopic>(this.url + "/" + topic);
  }
}
