package ng.lyu.pepperchat.comps

import ng.lyu.pepperchat.R
import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView

class ChatAdapter(val context: Context, val mData: ArrayList<data>) : RecyclerView.Adapter<ChatAdapter.ViewHolder>() {

    data class ViewHolder(
        val itemView: View,
        val sender: TextView,
        val chat1: TextView,
        val time1: TextView,
        val chat2: TextView,
        val time2: TextView,
        val my: View,
        val your: View
    ) : RecyclerView.ViewHolder(itemView)

    data class data(
        val sender: String,
        val chat: String,
        val time: String,
        val self: Boolean
    )

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = inflater.inflate(
            R.layout.chat_item,
            parent,
            false
        )
        val chat1: TextView = view.findViewById(R.id.text_chat1)
        val time1: TextView = view.findViewById(R.id.text_time1)
        val chat2: TextView = view.findViewById(R.id.text_chat2)
        val time2: TextView = view.findViewById(R.id.text_time2)
        val sender: TextView = view.findViewById(R.id.text_sender)
        val my: View = view.findViewById(R.id.my_chat)
        val your: View = view.findViewById(R.id.your_chat)
        return ViewHolder(view, sender, chat1, time1, chat2, time2, my, your)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        holder.sender.setText(mData.get(position).sender)
        holder.chat1.setText(mData.get(position).chat)
        holder.time1.setText(mData.get(position).time)
        holder.chat2.setText(mData.get(position).chat)
        holder.time2.setText(mData.get(position).time)

        if (mData.get(position).self) {
            holder.your.visibility = View.GONE
            holder.my.visibility = View.VISIBLE
        }
        else {
            holder.your.visibility = View.VISIBLE
            holder.my.visibility = View.GONE
        }
    }

    override fun getItemCount(): Int {
        return mData.size
    }
}