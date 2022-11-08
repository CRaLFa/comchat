import * as readline from 'readline/promises';
import { stdin as input, stdout as output } from 'process';
import { credentials, ClientDuplexStream } from '@grpc/grpc-js';
import { ChatMessage } from '../generated/command_chat_pb';
import { CommandChatClient } from '../generated/command_chat_grpc_pb';

const SERVER_ADDR = 'localhost:50051';

class CustomCommandChatClient extends CommandChatClient {
    user: string;

    constructor(userName: string) {
        super(SERVER_ADDR, credentials.createInsecure());
        this.user = userName;
    }

    async runChat() {
        const stream = this.chat();
        this.receive(stream);
        await this.sendMessage(stream, 'LOGGED_IN');
        await this.send(stream);
    }

    receive(stream: ClientDuplexStream<ChatMessage, ChatMessage>) {
        stream.on('data', (msg: ChatMessage) => {
            console.log('%s [%s] : %s', getDateTime(), msg.getAuthor(), msg.getBody());
        });
    }

    async send(stream: ClientDuplexStream<ChatMessage, ChatMessage>) {
        const rl = readline.createInterface({ input, output });

        for await (const line of rl) {
            const trimmed = line.trim();
            if (!trimmed) {
                continue;
            }
            if (trimmed === 'exit') {
                await this.sendMessage(stream, 'LOGGED_OUT');
                break;
            }
            await this.sendMessage(stream, trimmed);
        }

        rl.close();
    }

    async sendMessage(stream: ClientDuplexStream<ChatMessage, ChatMessage>, body: string) {
        await new Promise<void>((resolve, reject) => {
            const msg = new ChatMessage();
            msg.setAuthor(this.user);
            msg.setBody(body);

            stream.write(msg, (err: any) => {
                if (err) {
                    reject(err);
                } else {
                    resolve();
                }
            });
        });
    }
}

const getDateTime = () => {
    const padZero = (n: number, len: number) => {
        return ('0'.repeat(len) + n).slice(len * -1);
    };
    const now = new Date();

    const year = now.getFullYear();
    const month = padZero(now.getMonth() + 1, 2);
    const date = padZero(now.getDate(), 2);
    const hours = padZero(now.getHours(), 2);
    const minutes = padZero(now.getMinutes(), 2);
    const seconds = padZero(now.getSeconds(), 2);

    return `${year}/${month}/${date} ${hours}:${minutes}:${seconds}`;
}

const getUserName = async () => {
    const rl = readline.createInterface({ input, output });
    const name = await rl.question('Enter your name: ');
    rl.close();

    const trimmed = name.trim();
    return trimmed || '(anonymous)';
};

const main = async () => {
    const client = new CustomCommandChatClient(await getUserName());
    await client.runChat();

    client.close();
    process.exit(0);
};

main();
