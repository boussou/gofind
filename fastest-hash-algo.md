# Fastest hash of Non-Cryptographic algorithms


### **General Benchmark Results (Based on External Tests)**
| Algorithm  | Speed (GB/s) | Notes |
|------------|------------|-------|
| **xxHash** | ~29 GB/s  | Fastest, low CPU usage |
| **MurmurHash3** | ~15 GB/s | Good mix of speed & quality |
| **CityHash** | ~13 GB/s | Optimized for small strings |
| **FarmHash** | ~13-16 GB/s | Successor to CityHash, better for larger inputs |
| **MetroHash** | ~28 GB/s | Similar to xxHash, very fast |

### **Conclusions**
- **xxHash** is the fastest non-cryptographic hashing algorithm.
- **MetroHash** performs similarly but is less widely used.
- **MurmurHash** and **CityHash** offer a balance between speed and quality.


### **Theory**

The fastest cryptographic hash algorithms depend on use case (e.g., security, speed, and hardware support).

#### **1. For General-Purpose Fast Hashing (Non-Cryptographic)**
These are used for checksums, hash tables, and non-secure applications:
- **xxHash** – Extremely fast, optimized for speed.
- **MurmurHash** – High-performance hash, used in databases.
- **CityHash** – Developed by Google, optimized for short strings.
- **MetroHash** – Fast for both small and large inputs.
- **FarmHash** – Successor to CityHash, optimized for modern CPUs.

#### **2. For Cryptographic Fast Hashing (Secure)**
These are used for password hashing, data integrity, and cryptographic functions:
- **BLAKE3** – Fastest cryptographic hash function, parallelized, and highly secure.
- **SHA-3 (Keccak)** – Secure and efficient for cryptographic purposes.
- **SHA-256 (with hardware acceleration)** – Used in blockchain and security applications.
- **SipHash** – Fast for short input hashing, mainly for hash tables.

#### **3. Fastest Overall:**
- **BLAKE3** is the fastest cryptographic hash function due to its parallelization and efficiency.
- **xxHash** is the fastest non-cryptographic hash function, used for performance-sensitive applications.

For small inputs, **BLAKE3** overhead may be slightly higher. Otherwise It uses a tree-hash structure that can fully utilize multiple cores.

**Overall Trade-offs:**
  - If you need **extremely high-speed hashing** in contexts where security isn’t required (e.g., checksums, caching, deduplication), **xxHash** is often the better choice.
  - If you need a **cryptographically secure hash** that is still very fast—and especially if you can take advantage of multi-core parallelism—**BLAKE3** is an excellent option.

  
