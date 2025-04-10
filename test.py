import requests

def test_compress():
    url = "http://localhost:8080/compress"
    
    image_path = "sample.jpg" 
    
    with open(image_path, 'rb') as f:
        files = {'image': f}
        data = {'quality': '75', 'format': 'webp'}
        
        response = requests.post(url, files=files, data=data)
        
        print(f"Status Code: {response.status_code}")
        print(f"Response size: {len(response.content)} bytes")
        
        # Save compressed image
        with open("compressed.webp", "wb") as out_file:
            out_file.write(response.content)
            
if __name__ == "__main__":
    test_compress()