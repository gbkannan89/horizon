package com.example.mobile_flutter

import android.Manifest
import android.annotation.SuppressLint
import android.app.Activity
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.content.pm.PackageManager
import android.database.Cursor
import android.net.Uri
import android.provider.Telephony
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

class SmsPlugin(
    private val activity: Activity,
    private val engine: FlutterEngine
) {
    companion object {
        private const val CHANNEL = "com.thaari.horizon/sms"
        private const val SMS_PERMISSION_CODE = 1001
        private const val DAYS_BACK_DEFAULT = 30
    }

    private val channel: MethodChannel
    private var smsReceiver: SmsBroadcastReceiver? = null
    private var pendingCallback: MethodChannel.Result? = null

    // Keywords to filter bank transaction SMS
    private val transactionKeywords = listOf(
        "debited", "credited", "spent", "paid", "sent", "received",
        "upi", "neft", "imps", "trf", "transfer",
        "rs.", "rs ", "inr", "a/c", "account"
    )

    init {
        channel = MethodChannel(engine.dartExecutor.binaryMessenger, CHANNEL)
        channel.setMethodCallHandler { call, result ->
            when (call.method) {
                "hasSmsPermission" -> hasSmsPermission(result)
                "requestSmsPermission" -> requestSmsPermission(result)
                "readSmsInbox" -> {
                    val daysBack = call.argument<Int>("daysBack") ?: DAYS_BACK_DEFAULT
                    readSmsInbox(result, daysBack)
                }
                "startSmsListener" -> startSmsListener(result)
                "stopSmsListener" -> stopSmsListener(result)
                else -> result.notImplemented()
            }
        }
    }

    private fun hasSmsPermission(result: MethodChannel.Result) {
        val granted = ContextCompat.checkSelfPermission(
            activity, Manifest.permission.READ_SMS
        ) == PackageManager.PERMISSION_GRANTED
        result.success(granted)
    }

    private fun requestSmsPermission(result: MethodChannel.Result) {
        pendingCallback = result
        ActivityCompat.requestPermissions(
            activity,
            arrayOf(Manifest.permission.READ_SMS),
            SMS_PERMISSION_CODE
        )
    }

    fun onRequestPermissionsResult(requestCode: Int, permissions: Array<String>, grantResults: IntArray) {
        if (requestCode == SMS_PERMISSION_CODE) {
            val granted = grantResults.isNotEmpty() &&
                    grantResults[0] == PackageManager.PERMISSION_GRANTED
            pendingCallback?.success(granted)
            pendingCallback = null
        }
    }

    @SuppressLint("Range")
    private fun readSmsInbox(result: MethodChannel.Result, daysBack: Int) {
        if (ContextCompat.checkSelfPermission(activity, Manifest.permission.READ_SMS)
            != PackageManager.PERMISSION_GRANTED
        ) {
            result.success(emptyList<String>())
            return
        }

        try {
            val messages = mutableListOf<String>()
            val cutoffTime = System.currentTimeMillis() - (daysBack * 24L * 60 * 60 * 1000)

            val cursor: Cursor? = activity.contentResolver.query(
                Telephony.Sms.Inbox.CONTENT_URI,
                arrayOf(Telephony.Sms.Inbox.BODY, Telephony.Sms.Inbox.DATE),
                "${Telephony.Sms.Inbox.DATE} > ?",
                arrayOf(cutoffTime.toString()),
                "${Telephony.Sms.Inbox.DATE} DESC"
            )

            cursor?.use { c ->
                while (c.moveToNext()) {
                    val body = c.getString(c.getColumnIndex(Telephony.Sms.Inbox.BODY))
                    if (body != null && isTransactionSms(body)) {
                        messages.add(body)
                    }
                }
            }

            result.success(messages)
        } catch (e: Exception) {
            result.error("SMS_READ_ERROR", e.message, null)
        }
    }

    private fun isTransactionSms(body: String): Boolean {
        val lower = body.lowercase(Locale.ROOT)
        return transactionKeywords.any { lower.contains(it) }
    }

    private fun startSmsListener(result: MethodChannel.Result) {
        if (smsReceiver == null) {
            smsReceiver = SmsBroadcastReceiver(channel)
            val filter = IntentFilter(Telephony.Sms.Intents.SMS_RECEIVED_ACTION)
            // Use RECEIVE_SMS permission if available, otherwise register without
            try {
                activity.registerReceiver(smsReceiver, filter, RECEIVE_SMS_PERMISSION, null)
            } catch (e: Exception) {
                activity.registerReceiver(smsReceiver, filter)
            }
        }
        result.success(true)
    }

    private fun stopSmsListener(result: MethodChannel.Result) {
        smsReceiver?.let {
            try {
                activity.unregisterReceiver(it)
            } catch (_: Exception) {}
            smsReceiver = null
        }
        result.success(true)
    }

    fun destroy() {
        smsReceiver?.let {
            try {
                activity.unregisterReceiver(it)
            } catch (_: Exception) {}
        }
        smsReceiver = null
        channel.setMethodCallHandler(null)
    }

    private class SmsBroadcastReceiver(
        private val channel: MethodChannel
    ) : BroadcastReceiver() {
        override fun onReceive(context: Context, intent: Intent) {
            if (intent.action == Telephony.Sms.Intents.SMS_RECEIVED_ACTION) {
                val messages = Telephony.Sms.Intents.getMessagesFromIntent(intent)
                for (msg in messages) {
                    val body = msg.messageBody
                    if (body != null) {
                        channel.invokeMethod("onSmsReceived", body)
                    }
                }
            }
        }
    }
}

// Helper to get RECEIVE_SMS permission string safely
private val RECEIVE_SMS_PERMISSION: String? by lazy {
    try {
        Manifest.permission.RECEIVE_SMS
    } catch (_: Exception) {
        null
    }
}
