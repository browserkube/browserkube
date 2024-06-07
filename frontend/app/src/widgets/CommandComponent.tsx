import { Buffer } from 'buffer';
import styles from '@app/styles/commandsVideos.module.scss';
import { getMethodCommand } from '@shared/utils/getMethodCommand';

interface CommandComponentProps {
  method: string;
  command: string;
  request: string | null;
  statusCode: number;
  response: string;
}

const commandText = {
  fontSize: '13px',
};

const commandTitle = {
  fontWeight: '500',
  fontSize: '13px',
  marginBottom: '8px',
} as const;

export const CommandComponent = ({ props }: { props: CommandComponentProps }) => {
  const { method, command, request, statusCode, response } = props;
  return (
    <div className={styles.command_info}>
      <div style={{ ...commandTitle, marginBottom: '4px' }}> Command </div>
      <div className={styles.command_input}>{`${method} ${getMethodCommand(command)}`}</div>
      <div style={commandTitle}> Parametrs </div>
      <div className="command_text" style={{ ...commandText, marginBottom: '12px' }}>
        {request ? JSON.stringify(request, null, 2) : 'No avaliable data'}
      </div>
      <div style={commandTitle}>Response</div>
      <div className="command_text" style={commandText}>
        HTTPS Status: {statusCode}
      </div>
      <pre className="command_text" style={commandText}>
        {JSON.stringify(JSON.parse(Buffer.from(response, 'base64').toString('utf-8')), null, 2)}
      </pre>
    </div>
  );
};
