# upload
1. 获取图片列表 (/images/:order_id):
```swift
import Foundation

func getImages(orderID: String, completion: @escaping ([String]?, Error?) -> Void) {
    let urlString = "http://your-api-base-url/images/\(orderID)"
    guard let url = URL(string: urlString) else {
        completion(nil, NSError(domain: "Invalid URL", code: 0, userInfo: nil))
        return
    }
    
    URLSession.shared.dataTask(with: url) { (data, response, error) in
        if let error = error {
            completion(nil, error)
            return
        }
        
        guard let data = data else {
            completion(nil, NSError(domain: "No data received", code: 0, userInfo: nil))
            return
        }
        
        do {
            if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any],
               let images = json["images"] as? [String] {
                completion(images, nil)
            } else {
                completion(nil, NSError(domain: "Invalid JSON format", code: 0, userInfo: nil))
            }
        } catch {
            completion(nil, error)
        }
    }.resume()
}

// 使用示例
getImages(orderID: "123456") { (images, error) in
    if let error = error {
        print("Error: \(error.localizedDescription)")
    } else if let images = images {
        print("Images: \(images)")
    }
}

2. 上传图片 (/upload/:order_id):
```swift
import Foundation

func uploadImages(orderID: String, images: [UIImage], completion: @escaping (String?, Error?) -> Void) {
    let urlString = "http://your-api-base-url/upload/\(orderID)"
    guard let url = URL(string: urlString) else {
        completion(nil, NSError(domain: "Invalid URL", code: 0, userInfo: nil))
        return
    }
    
    let boundary = "Boundary-\(UUID().uuidString)"
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
    
    let httpBody = NSMutableData()
    
    for (index, image) in images.enumerated() {
        httpBody.append("--\(boundary)\r\n".data(using: .utf8)!)
        httpBody.append("Content-Disposition: form-data; name=\"images\"; filename=\"image_\(index + 1).jpg\"\r\n".data(using: .utf8)!)
        httpBody.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
        httpBody.append(image.jpegData(compressionQuality: 0.7)!)
        httpBody.append("\r\n".data(using: .utf8)!)
    }
    
    httpBody.append("--\(boundary)--\r\n".data(using: .utf8)!)
    request.httpBody = httpBody as Data
    
    URLSession.shared.dataTask(with: request) { (data, response, error) in
        if let error = error {
            completion(nil, error)
            return
        }
        
        guard let data = data else {
            completion(nil, NSError(domain: "No data received", code: 0, userInfo: nil))
            return
        }
        
        do {
            if let json = try JSONSerialization.jsonObject(with: data, options: []) as? [String: Any],
               let message = json["message"] as? String {
                completion(message, nil)
            } else {
                completion(nil, NSError(domain: "Invalid JSON format", code: 0, userInfo: nil))
            }
        } catch {
            completion(nil, error)
        }
    }.resume()
}

// 使用示例
let images = [UIImage(named: "image1")!, UIImage(named: "image2")!]
uploadImages(orderID: "123456", images: images) { (message, error) in
    if let error = error {
        print("Error: \(error.localizedDescription)")
    } else if let message = message {
        print("Message: \(message)")
    }
}
```

## android kotlin
1. 获取图片列表 (/images/:order_id):
```kotlin
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL

suspend fun getImages(orderID: String): List<String> = withContext(Dispatchers.IO) {
    val url = URL("http://your-api-base-url/images/$orderID")
    val connection = url.openConnection() as HttpURLConnection
    connection.requestMethod = "GET"
    
    try {
        val inputStream = connection.inputStream
        val response = inputStream.bufferedReader().use { it.readText() }
        val jsonObject = JSONObject(response)
        val imagesArray = jsonObject.getJSONArray("images")
        
        List(imagesArray.length()) { index ->
            imagesArray.getString(index)
        }
    } finally {
        connection.disconnect()
    }
}

// 使用示例
lifecycleScope.launch {
    try {
        val images = getImages("123456")
        println("Images: $images")
    } catch (e: Exception) {
        println("Error: ${e.message}")
    }
}
```

1. 上传图片 (/upload/:order_id):
```kotlin
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.json.JSONObject
import java.io.File
import java.net.HttpURLConnection
import java.net.URL

suspend fun uploadImages(orderID: String, images: List<File>): String = withContext(Dispatchers.IO) {
    val boundary = "Boundary-${System.currentTimeMillis()}"
    val url = URL("http://your-api-base-url/upload/$orderID")
    val connection = url.openConnection() as HttpURLConnection
    connection.requestMethod = "POST"
    connection.doOutput = true
    connection.setRequestProperty("Content-Type", "multipart/form-data; boundary=$boundary")
    
    try {
        connection.outputStream.use { outputStream ->
            images.forEachIndexed { index, image ->
                val fieldName = "images"
                val fileName = "image_${index + 1}.jpg"
                
                outputStream.write("--$boundary\r\n".toByteArray())
                outputStream.write("Content-Disposition: form-data; name=\"$fieldName\"; filename=\"$fileName\"\r\n".toByteArray())
                outputStream.write("Content-Type: image/jpeg\r\n\r\n".toByteArray())
                image.inputStream().use { it.copyTo(outputStream) }
                outputStream.write("\r\n".toByteArray())
            }
            outputStream.write("--$boundary--\r\n".toByteArray())
        }
        
        val inputStream = connection.inputStream
        val response = inputStream.bufferedReader().use { it.readText() }
        val jsonObject = JSONObject(response)
        jsonObject.getString("message")
    } finally {
        connection.disconnect()
    }
}

// 使用示例
lifecycleScope.launch {
    try {
        val images = listOf(File("/path/to/image1.jpg"), File("/path/to/image2.jpg"))
        val message = uploadImages("123456", images)
        println("Message: $message")
    } catch (e: Exception) {
        println("Error: ${e.message}")
    }
}
```