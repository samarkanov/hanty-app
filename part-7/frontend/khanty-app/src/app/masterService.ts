import { HttpClient, HttpResponse, HttpHeaders  } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Subscribers {
    topic: string
    subscribers: string[]
}

export interface SubscriptionForm {
    topic: string
    portno: string
}

export interface ActionEvent {
    success: boolean
    action: string
    client: string
    topic: string
    ip: string
}

export interface NotificationMessage {
    message: string
}

export class MasterService {
  private host = "https://develop.valenoq.com/sandbox/khanty-app/master"
  private url: string

  constructor(private http: HttpClient) {
      this.url = this.host
  }

  private generate_headers(): HttpHeaders {
      let headers: HttpHeaders = new HttpHeaders();
      headers = headers.append('Content-Type', 'application/x-www-form-urlencoded');
      headers = headers.append('X-CSRFToken', this.getCookie("csrftoken"));
      return headers
  }

  private getCookie(name: string) {
        let ca: Array<string> = document.cookie.split(';');
        let caLen: number = ca.length;
        let cookieName = `${name}=`;
        let c: string;

        for (let i: number = 0; i < caLen; i += 1) {
            c = ca[i].replace(/^\s+/g, '');
            if (c.indexOf(cookieName) == 0) {
                return c.substring(cookieName.length, c.length);
            }
        }
        return '';
    }

  public is_subscribed(topic:string) : Observable<Subscribers> {
      const url = this.url + "/subscribers/" + topic
      return this.http.get<Subscribers>(url);
  }

  public get_latest_notification() : Observable<NotificationMessage> {
      const url = this.url + "/get-notification"
      return this.http.get<NotificationMessage>(url);
  }

  public notify(topic:string, value:string) : void {
      const url = this.url + "/notify"
      let body = `topic=${topic}&value=${value}`
      const myheaders = {headers: this.generate_headers()}
      this.http.post(url, body, myheaders).subscribe(
          (res) => {},
          (err) => {console.log(err)}
      )
  }

  public sub(topic:string, portno:string) : Observable<ActionEvent> {
      const url = this.url + "/subscribe"
      let body = `topic=${topic}&portno=${portno}`
      const myheaders = {headers: this.generate_headers()}
      return this.http.post<ActionEvent>(url, body, myheaders);
  }

  public unsub(topic:string, portno:string) : Observable<ActionEvent> {
      const url = this.url + "/unsubscribe/" + topic + "/" + portno
      const myheaders = {headers: this.generate_headers()}
      return this.http.delete<ActionEvent>(url, myheaders);
  }

  public store_notiffication_msg(data: ActionEvent) : void {
      const url = this.url + "/store-notification"
      let body = `ip=${data.ip.split(":")[0]}&action=${data.action}&topic=${data.topic}&client=${data.client}`
      const myheaders = {headers: this.generate_headers()}
      // console.log("about to store notification, body: " + body)
      this.http.post(url, body, myheaders).subscribe(
          (res) => {},
          (err) => {console.log(err)}
      );
  }

}
