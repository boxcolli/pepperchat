package ng.lyu.pepperchat.comps

import android.app.Activity
import android.content.Intent
import android.util.Log
import ng.lyu.pepperchat.activities.ChatActivity
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString

class Socket(val nickname: String, val chat: String, var activity: Activity) : WebSocketListener() {
    override fun onOpen(webSocket: WebSocket, response: Response?) {
        socket = this
        Socket.websocket = websocket

        val intent = Intent(activity, ChatActivity::class.java)
        activity.startActivity(intent)
        activity.finish()
    }

    override fun onMessage(webSocket: WebSocket, bytes: ByteString) {
        val data = bytes.utf8()
        if (activity is ChatActivity) {
            (activity as ChatActivity).getMessage(data)
        }
    }

    override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
        webSocket.close(NORMAL_CLOSURE_STATUS, null)
        webSocket.cancel()
        activity?.finishAffinity()
    }

    override fun onFailure(webSocket: WebSocket?, t: Throwable, response: Response?) {
        activity?.finishAffinity()
    }

    companion object {
        const val NORMAL_CLOSURE_STATUS = 1000
        var socket: Socket? = null
        var websocket: WebSocket? = null
    }
}