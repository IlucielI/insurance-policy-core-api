#!/usr/bin/env python3
"""
Test semantic search endpoint
"""
import requests
import json

# API base URL
BASE_URL = "http://localhost:8080/api/v1"

def test_semantic_search():
    """Test the semantic search endpoint"""
    print("🔍 Testing semantic search endpoint...")
    
    # Test query
    query = "asuransi untuk keluarga"
    
    url = f"{BASE_URL}/products/search"
    payload = {
        "query": query,
        "limit": 5
    }
    
    print(f"\n📤 Request:")
    print(f"  URL: {url}")
    print(f"  Query: {query}")
    print(f"  Limit: 5")
    
    try:
        response = requests.post(url, json=payload, timeout=30)
        
        print(f"\n📥 Response:")
        print(f"  Status: {response.status_code}")
        
        if response.status_code == 200:
            data = response.json()
            print(f"  Results count: {data.get('count', 0)}")
            
            print(f"\n✅ Search Results:")
            for idx, product in enumerate(data.get('results', []), 1):
                print(f"\n  [{idx}] {product['name']}")
                print(f"      Category: {product['category']}")
                print(f"      Description: {product['description'][:100]}...")
                print(f"      ID: {product['id']}")
        else:
            print(f"  Error: {response.text}")
            
    except Exception as e:
        print(f"  ❌ Error: {e}")

if __name__ == "__main__":
    test_semantic_search()
