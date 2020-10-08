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
  private port = "10001"
  private host = "http://develop.valenoq.com"
  private url: string

  constructor(private http: HttpClient) {
      this.url = this.host + ":" + this.port
  }

  public listTopics() : Observable<ListTopics> {
      return this.http.get<ListTopics>(this.url);
  }

  public statesForTopic(topic: string) : Observable<StateTopic> {
      return this.http.get<StateTopic>(this.url + "/" + topic);
      // var url = {
      //     topic: ajax.getJSON(this.url + "/" + topic)
      // }
      // return this.http.get<StateTopic>(url).pipe(
      //     mergeMap(x => interval(1000).pipe(map(i => x)))
      // )
      // var urls = {
      //     url:
      // }
      // for (const topicName of topics) {
      //     urls[topic] = ajax.getJSON(url + "/" + topicName)
      // }
      // return interval(1000).pipe(switchMap(() => forkJoin(url)))
  }
}
