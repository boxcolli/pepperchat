package ng.lyu.pepperchat.activities

import android.annotation.SuppressLint
import android.os.Bundle
import android.text.InputType
import android.util.Log
import android.view.KeyEvent
import android.view.View
import android.view.View.OnKeyListener
import android.view.inputmethod.InputMethodManager
import android.widget.EditText
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import ng.lyu.pepperchat.R
import ng.lyu.pepperchat.comps.ChatAdapter
import ng.lyu.pepperchat.comps.Socket
import okhttp3.MediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import org.json.JSONObject
import java.io.BufferedReader
import java.io.DataOutputStream
import java.io.InputStreamReader
import java.net.HttpURLConnection
import java.net.URL
import kotlin.concurrent.thread


class ChatActivity(
    private var chatList: RecyclerView? = null,
    private var mData: ArrayList<ChatAdapter.data>? = null,
    private var adapter: ChatAdapter? = null,
    private var layoutManager: LinearLayoutManager? = null,

    private var inputChat: EditText? = null
) : AppCompatActivity(), OnKeyListener {

    // activity functions
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_chat)
        Socket.socket?.activity = this

        initChatList()
        initText()
    }

    // events
    override fun onKey(p0: View?, p1: Int, p2: KeyEvent?): Boolean {
        if (p2?.keyCode == KeyEvent.KEYCODE_ENTER && p2.action == KeyEvent.ACTION_DOWN) {
            val text = inputChat?.text?.toString() ?: ""
            thread(start = true) {
                try {
                    val client = OkHttpClient()
                    val req = Request.Builder()
                        .url("http://10.0.2.2:8080/chat/${Socket.socket?.chat}/message")
                        .addHeader("x-user-id", Socket.socket?.nickname ?: "")
                        .post(RequestBody.create(MediaType.parse("application/json"), "{\"content\":\"${text}\"}"))
                        .build()
                    client.newCall(req).execute()
                }
                catch (e: java.lang.Exception) {
                    Log.e("Socket", e.toString())
                }
            }

            inputChat?.text?.clear()
            return true;
        }
        return false;
    }

    override fun finish() {
        super.finish()
        Socket.websocket?.close(Socket.NORMAL_CLOSURE_STATUS, null)
        Socket.websocket?.cancel()
    }

    // chat activity functions
    @SuppressLint("NotifyDataSetChanged")
    private fun appendChat(data: ChatAdapter.data) {
        val bottom = (chatList?.layoutManager as LinearLayoutManager).findLastCompletelyVisibleItemPosition()
        mData?.let { it ->
            it.add(data)
            adapter!!.notifyDataSetChanged()

            if (bottom == it.size.minus(2)) {
                chatList?.scrollToPosition(it.size - 1)
            }
        }
    }

    private fun initText() {
        inputChat = findViewById(R.id.input_chat)
        inputChat?.setOnKeyListener(this)
        inputChat?.inputType = InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_FLAG_MULTI_LINE

        (findViewById(R.id.text_room) as TextView).text = Socket.socket?.chat
    }

    private fun initChatList() {
        chatList = findViewById(R.id.chats)
        mData = arrayListOf()
        adapter = ChatAdapter(this, mData!!)
        layoutManager = LinearLayoutManager(this)
        chatList?.layoutManager = layoutManager
        chatList?.adapter = adapter
        chatList?.scrollToPosition(mData!!.size - 1)
    }

    fun getMessage(data: String) {
        try {
            val json = JSONObject(data)
            val sender = json.get("username").toString()
            val content = json.get("content").toString()
            val time = json.get("created_time").toString().split(" ")[0]

            runOnUiThread {
                appendChat(ChatAdapter.data(sender, content, time, sender == Socket.socket?.nickname))
            }
        }
        catch (e: java.lang.Exception) {
            e.printStackTrace()
        }
    }
}