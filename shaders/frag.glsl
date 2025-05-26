#version 460
out vec4 frag_color;

const int MAX_RAY_STEPS = 128;
const int WORLD_SIZE = 512;
uniform vec2 uSize;
uniform float uTime;
uniform vec3 uPlayerPos;
const float PI = 3.14159265359;
uniform mat4 uInvView;
uniform mat4 uInvProj;

struct FlatNode {
    int child_indices[8]; // -1 means no child
    int is_leaf; // 1 for leaf, 0 otherwise
    int _pad0, _pad1, _pad2; // pad to mutliple of 16 bytes (std430)
    // optional: int data; for voxel value
};

layout(std430, binding = 0) buffer OctreeBuffer {
    FlatNode nodes[];
};

// Octree ray traversal algorithm
// Returns bvec3 mask indicating which face was hit (x=1 means hit X face, etc.)

struct TraversalState {
    int node_index;
    vec3 node_min;
    vec3 node_size;
    float t_min;
    float t_max;
};

// Stack for traversal (adjust size based on max octree depth)
const int MAX_STACK_SIZE = 32;
TraversalState stack[MAX_STACK_SIZE];
int stack_ptr = 0;

// Ray-AABB intersection test
bool ray_aabb_intersect(vec3 ray_origin, vec3 ray_dir, vec3 box_min, vec3 box_max, out float t_near, out float t_far) {
    vec3 inv_dir = 1.0 / ray_dir;
    vec3 t1 = (box_min - ray_origin) * inv_dir;
    vec3 t2 = (box_max - ray_origin) * inv_dir;

    vec3 t_min_vec = min(t1, t2);
    vec3 t_max_vec = max(t1, t2);

    t_near = max(max(t_min_vec.x, t_min_vec.y), t_min_vec.z);
    t_far = min(min(t_max_vec.x, t_max_vec.y), t_max_vec.z);

    return t_near <= t_far && t_far > 0.0;
}

// Get which face was hit based on ray intersection
bvec3 get_hit_face(vec3 ray_origin, vec3 ray_dir, vec3 box_min, vec3 box_max, float t_hit) {
    vec3 hit_point = ray_origin + ray_dir * t_hit;
    vec3 center = (box_min + box_max) * 0.5;
    vec3 size = box_max - box_min;

    // Normalize hit point relative to box center
    vec3 rel_hit = (hit_point - center) / (size * 0.5);

    // Find which component has the largest absolute value
    vec3 abs_rel = abs(rel_hit);
    float max_component = max(max(abs_rel.x, abs_rel.y), abs_rel.z);

    // Tolerance for floating point comparison
    const float epsilon = 1e-6;

    return bvec3(
        abs(abs_rel.x - max_component) < epsilon,
        abs(abs_rel.y - max_component) < epsilon,
        abs(abs_rel.z - max_component) < epsilon
    );
}

// Get child index for octree traversal
int get_child_index(vec3 point, vec3 node_center) {
    int index = 0;
    if (point.x >= node_center.x) index |= 1;
    if (point.y >= node_center.y) index |= 2;
    if (point.z >= node_center.z) index |= 4;
    return index;
}

// Main octree traversal function
bvec3 traverse_octree(vec3 ray_origin, vec3 ray_direction) {
    // Normalize ray direction
    vec3 ray_dir = normalize(ray_direction);

    // Initialize traversal with root node
    vec3 world_min = vec3(0.0);
    vec3 world_max = vec3(WORLD_SIZE);

    float t_near, t_far;
    if (!ray_aabb_intersect(ray_origin, ray_dir, world_min, world_max, t_near, t_far)) {
        return bvec3(false); // No intersection with octree bounds
    }

    // Initialize stack with root
    stack_ptr = 0;
    stack[stack_ptr].node_index = 0;
    stack[stack_ptr].node_min = world_min;
    stack[stack_ptr].node_size = vec3(WORLD_SIZE);
    stack[stack_ptr].t_min = max(t_near, 0.0);
    stack[stack_ptr].t_max = t_far;
    stack_ptr++;

    float closest_t = 1e30;
    bvec3 hit_face = bvec3(false);

    // Traversal loop
    while (stack_ptr > 0) {
        stack_ptr--;
        TraversalState current = stack[stack_ptr];

        // Skip if this node is further than our closest hit
        if (current.t_min >= closest_t) {
            continue;
        }

        FlatNode node = nodes[current.node_index];

        if (node.is_leaf == 1) {
            // Leaf node - we have a hit
            if (current.t_min < closest_t) {
                closest_t = current.t_min;
                vec3 box_max = current.node_min + current.node_size;
                hit_face = get_hit_face(ray_origin, ray_dir, current.node_min, box_max, current.t_min);
            }
        } else {
            // Internal node - traverse children
            vec3 node_center = current.node_min + current.node_size * 0.5;
            vec3 child_size = current.node_size * 0.5;

            // Order children by ray direction for front-to-back traversal
            int child_order[8];
            for (int i = 0; i < 8; i++) {
                child_order[i] = i;
            }

            // Simple ordering based on ray direction signs
            if (ray_dir.x < 0.0) {
                for (int i = 0; i < 8; i++) {
                    child_order[i] ^= 1;
                }
            }
            if (ray_dir.y < 0.0) {
                for (int i = 0; i < 8; i++) {
                    child_order[i] ^= 2;
                }
            }
            if (ray_dir.z < 0.0) {
                for (int i = 0; i < 8; i++) {
                    child_order[i] ^= 4;
                }
            }

            // Add children to stack in reverse order for proper traversal
            for (int i = 7; i >= 0; i--) {
                int child_idx = child_order[i];
                int child_node_index = node.child_indices[child_idx];

                if (child_node_index != -1 && stack_ptr < MAX_STACK_SIZE) {
                    // Calculate child bounds
                    vec3 child_min = current.node_min;
                    if ((child_idx & 1) != 0) child_min.x += child_size.x;
                    if ((child_idx & 2) != 0) child_min.y += child_size.y;
                    if ((child_idx & 4) != 0) child_min.z += child_size.z;
                    vec3 child_max = child_min + child_size;

                    // Test ray intersection with child
                    float child_t_near, child_t_far;
                    if (ray_aabb_intersect(ray_origin, ray_dir, child_min, child_max, child_t_near, child_t_far)) {
                        if (child_t_near < closest_t) {
                            stack[stack_ptr].node_index = child_node_index;
                            stack[stack_ptr].node_min = child_min;
                            stack[stack_ptr].node_size = child_size;
                            stack[stack_ptr].t_min = max(child_t_near, 0.0);
                            stack[stack_ptr].t_max = child_t_far;
                            stack_ptr++;
                        }
                    }
                }
            }
        }
    }

    return hit_face;
}
void main() {
    vec2 uv = (gl_FragCoord.xy * 2.0 - uSize) / uSize.y;
    vec3 rayOrigin = uPlayerPos;

    vec4 target = uInvProj * vec4(uv, 1, 1);
    vec3 rayDirection = (uInvView * vec4(normalize(target.xyz / target.w), 0)).xyz;
    vec3 lightDirection = -1 * normalize(vec3(sin(uTime * 2), 1, cos(uTime * 2)));

    // frag_color = vec4(raySphereIntersection(rayOrigin, rayDirection, .7, lightDirection), 1);
    // vec3 texCoord = vec3((uv.x + 1.0) / 2.0, (uv.y + 1.0) / 2.0, 0);
    // frag_color = vec4(vec3(texture(voxelMap, vec3(uv * 4, 0)) * 255), 1);
    bvec3 mask = traverse_octree(rayOrigin, rayDirection);
    if (mask.x) {
        frag_color.xyz = vec3(0.5);
    }
    else if (mask.y) {
        frag_color.xyz = vec3(1.0);
    }
    else if (mask.z) {
        frag_color.xyz = vec3(0.75);
    }
    else {
        frag_color.xyz = vec3(0.53, 0.81, 0.94);
    }
    frag_color.w = 1;
}
