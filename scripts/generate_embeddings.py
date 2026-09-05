#!/usr/bin/env python3
"""
Generate embeddings for all products in the database
"""
import os
import sys
import json
import requests
import psycopg2
from psycopg2.extras import execute_values

def get_db_connection():
    """Connect to the database"""
    # Read from environment variable
    db_url = os.getenv("DATABASE_URL")
    if not db_url:
        print("❌ DATABASE_URL environment variable not set")
        sys.exit(1)
    
    # Parse connection string
    db_url = db_url.replace("postgres://", "postgresql://")
    return psycopg2.connect(db_url)

def generate_embedding(text, base_url, model):
    """Generate embedding using LLM API"""
    url = f"{base_url}/embeddings"
    payload = {
        "input": text,
        "model": model
    }
    
    response = requests.post(url, json=payload, timeout=30)
    response.raise_for_status()
    
    data = response.json()
    if not data.get("data") or len(data["data"]) == 0:
        raise Exception("No embedding data in response")
    
    return data["data"][0]["embedding"]

def main():
    print("🚀 Starting embeddings generation for products...")
    
    # Configuration
    llm_base_url = os.getenv("LLM_BASE_URL", "http://100.103.220.104:20128/v1")
    embeddings_model = os.getenv("EMBEDDINGS_MODEL", "bge-m3")
    
    print(f"🔍 Embeddings Config: {llm_base_url} (model: {embeddings_model})")
    
    # Connect to database
    try:
        conn = get_db_connection()
        cur = conn.cursor()
        print("✅ Database connected")
    except Exception as e:
        print(f"❌ Database connection failed: {e}")
        sys.exit(1)
    
    # Add unique constraint if it doesn't exist
    try:
        cur.execute("""
            ALTER TABLE product_embeddings 
            ADD CONSTRAINT unique_product_chunk_type UNIQUE (product_id, chunk_type)
        """)
        conn.commit()
        print("✅ Added unique constraint to product_embeddings")
    except Exception as e:
        conn.rollback()
        if "already exists" in str(e):
            print("⚠️  Unique constraint already exists")
        else:
            print(f"⚠️  Could not add constraint: {e}")
    
    # Fetch all products
    cur.execute("""
        SELECT id, name, slug, category, description, coverage_details
        FROM products
        ORDER BY created_at DESC
    """)
    
    products = cur.fetchall()
    print(f"📦 Found {len(products)} products")
    
    if len(products) == 0:
        print("⚠️  No products found in database")
        cur.close()
        conn.close()
        return
    
    # Generate embeddings for each product
    success_count = 0
    for idx, (product_id, name, slug, category, description, coverage_details) in enumerate(products, 1):
        print(f"\n[{idx}/{len(products)}] Processing: {name}")
        
        # Create searchable text
        search_text = f"{name}. {description}"
        if category:
            search_text = f"Kategori: {category}. {search_text}"
        
        print(f"  Text: {search_text[:100]}...")
        
        try:
            # Generate embedding
            embedding = generate_embedding(search_text, llm_base_url, embeddings_model)
            print(f"  ✅ Generated embedding (dim: {len(embedding)})")
            
            # Convert embedding to PostgreSQL array format
            embedding_str = "[" + ",".join(str(x) for x in embedding) + "]"
            
            # Save to database
            cur.execute("""
                INSERT INTO product_embeddings (product_id, chunk_type, chunk_text, embedding)
                VALUES (%s, %s, %s, %s::vector)
                ON CONFLICT (product_id, chunk_type) 
                DO UPDATE SET chunk_text = %s, embedding = %s::vector, created_at = NOW()
            """, (product_id, "description", search_text, embedding_str, search_text, embedding_str))
            
            conn.commit()
            print(f"  ✅ Saved embedding to database")
            success_count += 1
            
        except Exception as e:
            print(f"  ❌ Error: {e}")
            conn.rollback()
            continue
    
    print(f"\n✅ Embeddings generation completed! ({success_count}/{len(products)} successful)")
    
    # Verify embeddings
    cur.execute("SELECT COUNT(*) FROM product_embeddings")
    result = cur.fetchone()
    count = result[0] if result else 0
    print(f"📊 Total embeddings in database: {count}")
    
    cur.close()
    conn.close()

if __name__ == "__main__":
    main()
