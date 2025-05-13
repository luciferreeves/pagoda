export interface District {
  id: string;
  name: string;
  description: string;
  image: string;
}

export interface Site {
  name: string;
  url: string;
  description: string;
  owner: string;
  tags: string[];
  added: Date;
  screenshotUrl: string;
}
