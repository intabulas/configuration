package hocon

import (
	"bytes"
	"fmt"
	"strings"
)

type HoconObject struct {
	items map[string]*HoconValue
	keys  []string
}

func NewHoconObject() *HoconObject {
	return &HoconObject{
		items: make(map[string]*HoconValue),
	}
}

func (p *HoconObject) GetString() string {
	panic("This element is an object and not a string.")
}

func (p *HoconObject) IsArray() bool {
	return false
}

func (p *HoconObject) GetArray() []*HoconValue {
	panic("This element is an object and not an array.")
}

func (p *HoconObject) GetKeys() []string {
	return p.keys
}

func (p *HoconObject) Unwrapped() map[string]interface{} {
	if len(p.items) == 0 {
		return nil
	}

	dics := map[string]interface{}{}

	for _, k := range p.keys {
		v := p.items[k]

		obj := v.GetObject()
		if obj != nil {
			dics[k] = obj.Unwrapped()
		} else {
			dics[k] = v
		}
	}

	return dics
}

func (p *HoconObject) Items() map[string]*HoconValue {
	return p.items
}

func (p *HoconObject) GetKey(key string) *HoconValue {
	value, _ := p.items[key]
	return value
}

func (p *HoconObject) GetOrCreateKey(key string) *HoconValue {
	if value, exist := p.items[key]; exist {
		child := NewHoconValue()
		child.oldValue = value
		p.items[key] = child
		return child
	}

	child := NewHoconValue()
	p.items[key] = child
	p.keys = append(p.keys, key)
	return child
}

func (p *HoconObject) IsString() bool {
	return false
}

func (p *HoconObject) String() string {
	return p.ToString(0)
}

func (p *HoconObject) ToString(indent int) string {
	tmp := strings.Repeat(" ", indent*2)
	buf := bytes.NewBuffer(nil)
	for _, k := range p.keys {
		key := p.quoteIfNeeded(k)
		v := p.items[key]
		buf.WriteString(fmt.Sprintf("%s%s : %s\r\n", tmp, key, v.ToString(indent)))
	}
	return buf.String()
}

func (p *HoconObject) Merge(other *HoconObject) {
	thisValues := p.items
	otherItems := other.items

	otherKeys := other.keys

	for _, otherkey := range otherKeys {

		otherValue := otherItems[otherkey]
		if thisValue, exist := thisValues[otherkey]; exist {
			if thisValue.IsObject() && otherValue.IsObject() {
				thisValue.GetObject().Merge(otherValue.GetObject())
			}
		} else {
			p.items[otherkey] = otherValue
			p.keys = append(p.keys, otherkey)
		}
	}
}

func (p *HoconObject) MergeImmutable(other *HoconObject) *HoconObject {
	thisValues := map[string]*HoconValue{}
	otherKeys := other.keys

	var thisKeys []string

	otherItems := other.items

	for _, otherkey := range otherKeys {
		otherValue := otherItems[otherkey]

		if thisValue, exist := thisValues[otherkey]; exist {

			if thisValue.IsObject() && otherValue.IsObject() {

				mergedObject := thisValue.GetObject().MergeImmutable(otherValue.GetObject())
				mergedValue := NewHoconValue()

				mergedValue.AppendValue(mergedObject)
				thisValues[otherkey] = mergedValue
			}
		} else {
			thisValues[otherkey] = &HoconValue{values: otherValue.values}
			thisKeys = append(thisKeys, otherkey)
		}
	}

	return &HoconObject{items: thisValues, keys: thisKeys}
}

func (p *HoconObject) quoteIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
}
