package ng.lyu.pepperchat.activities

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.view.View.OnClickListener
import android.widget.EditText
import android.widget.TextView
import android.widget.Toast
import ng.lyu.pepperchat.R
import ng.lyu.pepperchat.comps.Socket
import okhttp3.OkHttpClient
import okhttp3.Request

class LoginActivity(
    private var input_nickname: EditText? = null,
    private var input_room: EditText? = null,
    private var button_join: TextView? = null
) : AppCompatActivity(), OnClickListener {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_login)

        input_nickname = findViewById(R.id.input_nickname)
        input_room = findViewById(R.id.input_room)
        button_join = findViewById(R.id.button_join)

        button_join?.setOnClickListener(this)
    }

    override fun onClick(p0: View?) {
        if (input_nickname?.isEnabled == false) {
            return
        }
        if (input_nickname?.text?.isEmpty() == true) {
            Toast.makeText(this, "Please enter your nickname.", Toast.LENGTH_SHORT).show()
            return
        }
        if (input_room?.text?.isEmpty() == true) {
            Toast.makeText(this, "Please enter room name.", Toast.LENGTH_SHORT).show()
            return
        }
        enterRoom(input_nickname!!.text.toString(), input_room!!.text.toString())
    }

    override fun finish() {
        super.finish()
        Socket.websocket?.close(Socket.NORMAL_CLOSURE_STATUS, null)
        Socket.websocket?.cancel()
    }

    private fun enterRoom(nickname: String, room: String) {
        input_nickname?.isEnabled = false
        input_room?.isEnabled = false

        val client = OkHttpClient()
        val request = Request.Builder()
            .url("ws://10.0.2.2:8080/chat/${room}/ws")
            .build()
        val listener = Socket(nickname, room, this)

        client.newWebSocket(request, listener)
        client.dispatcher().executorService().shutdown()
    }
}