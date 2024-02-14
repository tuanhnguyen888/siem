export enum ICurrentPage {
  Event = "event",
  Alert = "alert",
}

export interface IEvent {
  id?: string | number;
  data?: object | any;
  localTime?: string | number;
  source?: string;
  utcTime?: string | number;
}

export interface IAlert {
  id?: string | number;
  eventId?: string | number;
  data?: object | any;
  level?: string | number;
  message?: string;
  localTime?: string | number;
}

export interface IUser {
  id?: number | string;
  username?: string;
  password?: string;
  fullname?: string;
  avarar?: string;
  userRole?: number;
}

export interface ILogSource {
  id?: number | string;
  inputName?: string;
  properties?: object | any;
  createdAt?: number;
  updatedAt?: number;
  protocol?: string;
}
