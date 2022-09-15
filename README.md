# TCP-Chat

Motivação:
A segunda parte do trabalho consiste em implementar um Chat usando TCP. O Chat deve suportar múltiplos clientes e um servidor. Todos os clientes devem estar na mesma sala do chat (i.e., as mensagens enviadas por um cliente devem ser recebidas por todos os clientes). Comandos que o usuário (i.e., cliente) pode enviar:
<ul>
  <li>
    /ENTRAR: ao usar esse comando, é requisitado o IP e porta do servidor, além do nickname que o usuário deseja usar no chat (não precisa tratar nicknames repetidos). Todos os usuários devem ser notificados da entrada do novo usuário;
  </li>
   <li>
     Uma vez conectado, o cliente pode enviar mensagens para a sala do chat e deve receber todas as mensagens que forem enviadas pelos outros usuários;
  </li>
   <li>
     /USUARIOS: ao enviar esse comando, o cliente recebe a lista de usuários atualmente conectados ao chat;
  </li>
   <li>
     /SAIR: ao enviar esse comando, uma mensagem é enviada à sala do chat informando que o usuário está saindo e encerra a participação no chat.
  </li>
</ul>

É papel do servidor receber as requisições dos clientes e encaminhar as mensagens
recebidas para todos eles. Descreva o formato para cada tipo de mensagem. Não pode
usar comunicação em grupo
