export interface LetterParticipant {
  username: string;
  display_name: string;
  avatar_url: string;
  role: string;
}

export interface LetterAttachment {
  ref: string;
  file_name: string;
  url: string;
  file_size: number;
  content_type: string;
  category: string;
}

export interface LetterMessage {
  ref: string;
  sender: LetterParticipant | null;
  body: string;
  attachments: LetterAttachment[];
  edited_at: string | null;
  created_at: string;
  deleted: boolean;
}

export interface Letter {
  ref: string;
  title: string;
  is_system: boolean;
  system_ref?: string;
  participants: LetterParticipant[];
  last_message?: LetterMessage;
  unread: boolean;
  updated_at: string;
}

export interface LetterDetail {
  ref: string;
  title: string;
  is_system: boolean;
  system_ref?: string;
  participants: LetterParticipant[];
  messages: LetterMessage[];
  created_at: string;
}