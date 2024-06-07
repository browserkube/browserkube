package io.browserkube;

import org.java_websocket.handshake.ServerHandshake;

import java.net.URI;

public class WebSocketClient extends org.java_websocket.client.WebSocketClient {
    // Simple message buffer for testing: NOT THREAD SAFE!
    String message;
    public String getMessage() {
        return message;
    }
    public WebSocketClient(URI serverUri) {
        super(serverUri);
        this.message = new String();
    }
    @Override
    public void onOpen(ServerHandshake serverHandshake) {
        System.out.println("conn open");
    }
    @Override
    public void onMessage(String s) {
        System.out.println("msg > :" + s);
        this.message = s;
    }
    @Override
    public void onClose(int code, String reason, boolean remote) {
        System.out.println(
                "Connection closed by " + (remote ? "remote peer" : "us") + " Code: " + code + " Reason: "
                        + reason);
    }
    @Override
    public void onError(Exception e) {
        e.printStackTrace();
    }
}
