import { Component } from '@angular/core';
import { TopicService, ListTopics, StateTopic } from './topicService';
import { Config, ConfigItem, ConfigService } from './configService';
import { ClientData, ClientService, ClientTopicInfo } from './clientService';
import { MasterService, Subscribers, ActionEvent, NotificationMessage } from './masterService';
import { HttpClient, HttpResponse } from '@angular/common/http';

import { Observable, Subject, throwError, of, interval } from 'rxjs';
import { catchError, retry, debounceTime, distinctUntilChanged, switchMap, mergeMap, map } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {
  public title = 'khanty-app';
  public topics: ListTopics
  public currentState: Record<string, string> = {}
  public config: Config
  public clients: Record<string, ConfigItem> = {}
  public clNamePortMap: Record<string, string> = {}
  public latestNotification: string = ""

  private topicService : TopicService
  private configService: ConfigService
  private clientService: ClientService
  private masterService: MasterService

  constructor(private http: HttpClient){
      this.topicService = new TopicService(this.http)
      this.configService = new ConfigService(this.http)
      this.clientService = new ClientService(this.http)
      this.masterService = new MasterService(this.http)

      this.topicService.listTopics().subscribe(
          (data: ListTopics) => {
              this.topics = data
          },
          (err) => {console.log(err)},
          () => {this.initConfig()}
      )

      // notifications
      interval(800).subscribe(() => this.notify("ChangeColor"))
      interval(1600).subscribe(() => this.notify("SendMessage"))
  }

  private getClientsState() {
      for (const portno in this.clients) {

          // update message and color
          this.clientService.getState(portno, this.topics.topics).subscribe(
                // here we deal with an endless loop
                (data: ClientTopicInfo) => {
                     this.updateClientState(data)

                     // message checkbox update
                     this.masterService.is_subscribed("SendMessage").subscribe(
                         (data:Subscribers) => {
                             if (data.subscribers && data.subscribers.indexOf(portno) > -1) {
                                 this.clients[portno].messagesSubscribe = true
                             } else{
                                 this.clients[portno].messagesSubscribe = false
                             }
                     })

                     // color checkbox update
                     this.masterService.is_subscribed("ChangeColor").subscribe(
                         (data:Subscribers) => {
                             if (data.subscribers && data.subscribers.indexOf(portno) > -1) {
                                 this.clients[portno].colorSubscribe = true
                             } else{
                                 this.clients[portno].colorSubscribe = false
                             }
                     })

                     // notification pop-up
                     this.masterService.get_latest_notification().subscribe(
                         (data:NotificationMessage) => {
                             // console.log(data)
                             if (data) this.latestNotification = data.message;
                         }
                     )
          });


      }
  }

  private updateClientState(data: ClientTopicInfo) {
      const clientName = data["ChangeColor"]["name"]
      const clientPort = this.clNamePortMap[clientName]
      this.clients[clientPort].message = data["SendMessage"].data
      this.clients[clientPort].color = data["ChangeColor"].data
  }

  private initConfig() {
      this.configService.getConfig().subscribe(
          (data: Config) => {
              this.config = data
          },
          (err) => {console.log(err)},
          () => {this.initClients()}
      )
  }

  private initClients() {

      for (var serviceName in this.config) {
            if (serviceName.indexOf("client") > -1) {
                const portno = this.config[serviceName].port
                var item : ConfigItem = {
                    host: this.config[serviceName].host,
                    port: portno,
                    name: serviceName,
                    messagesSubscribe: false,
                    colorSubscribe: false
                }
                this.clients[item.port] = item
                this.clNamePortMap[item.name] = item.port
            }
      }
      this.getClientsState()
  }

  public onMsgEvent(portno: string, event: any){
      let reply = {}
      if (event.target.checked){
          this.masterService.sub("SendMessage", portno).subscribe(
              (data: ActionEvent) => {
                  data.action = "subscribed"
                  this.store_msg(data)
              }
          )
      } else {
          this.masterService.unsub("SendMessage", portno).subscribe(
              (data: ActionEvent) => {
                  data.action = "unsubscribed"
                  this.store_msg(data)
              }
          )
      }
  }

  private store_msg(data: ActionEvent) {
        this.masterService.store_notiffication_msg(data)
  }

  public onColorEvent(portno: string, event: any){
      if (event.target.checked){
          this.masterService.sub("ChangeColor", portno).subscribe(
              (data: ActionEvent) => {
                  data.action = "subscribed"
                  this.store_msg(data)
              }
          )
      } else {
          this.masterService.unsub("ChangeColor", portno).subscribe(
              (data: ActionEvent) => {
                  data.action = "unsubscribed"
                  this.store_msg(data)
              }
          )
      }
  }

  public notify(topic){
      // get non-repeatable random value out of all possible topic values
      this.topicService.statesForTopic(topic).subscribe(
          (data: StateTopic) => {
                  const idx = data.states.indexOf(this.currentState[topic])
                  if (idx > -1) {
                      data.states.splice(idx, 1)
                  }
                  const random = Math.floor(Math.random() * data.states.length)
                  this.currentState[topic] = data.states[random]
                  this.masterService.notify(topic, this.currentState[topic])
              }
      );
  }

}
