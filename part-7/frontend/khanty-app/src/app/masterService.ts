import { HttpClient, HttpResponse, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Subscribers {
    topic: string
    subscribers: string[]
}

export interface SubscriptionForm {
    topic: string
    portno: string
}

const HTTP_HEADERS = {
    headers: new HttpHeaders({
        'Content-Type': 'application/x-www-form-urlencoded'
    })
}

export class MasterService {
  private port = "10004"
  private host = "http://develop.valenoq.com"
  private url: string

  constructor(private http: HttpClient) {
      this.url = this.host + ":" + this.port
  }

  public is_subscribed(topic:string) : Observable<Subscribers> {
      const url = this.url + "/subscribers/" + topic
      return this.http.get<Subscribers>(url);
  }

  public notify(topic:string, value:string) : void {
      const url = this.url + "/notify"
      let body = `topic=${topic}&value=${value}`
      this.http.post(url, body, HTTP_HEADERS).subscribe(
          (res) => {},
          (err) => {console.log(err)}
      )
  }

  public subscribe(topic:string, portno:string) : void {
      const url = this.url + "/subscribe"
      let body = `topic=${topic}&portno=${portno}`
      this.http.post(url, body, HTTP_HEADERS).subscribe(
          (res) => {},
          (err) => {console.log(err)}
      )
  }

  public unsubscribe(topic:string, portno:string) : void {
      const url = this.url + "/unsubscribe/" + topic + "/" + portno

      this.http.delete(url).subscribe(
          (res) => {},
          (err) => {console.log(err)}
      )
  }
}
